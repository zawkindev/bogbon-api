{ pkgs }:

{
  network.description = "Django REST API with HTTPS and PostgreSQL";

  apiServer = { config, pkgs, ... }: {

    # --- Python Environment ---
    environment.systemPackages = with pkgs; [
      postgresql
      python312
      gunicorn
      python312Packages.django
      python312Packages.psycopg2
      python312Packages.djangorestframework
    ];

    # --- PostgreSQL Database ---
    services.postgresql = {
      enable = true;
      ensureDatabases = [ "bogbon" ];
      ensureUsers = [{
        name = "django-user";
        ensurePermissions = { "DATABASE bogbon" = "ALL PRIVILEGES"; };
      }];
      authentication = ''
        # Allow password auth for Django user
        local all django-user md5
      '';
      initialScript = pkgs.writeText "postgres-init.sql" ''
        ALTER USER "django-user" WITH PASSWORD '1';
      '';
    };


    # --- Backup ---
    services.postgresqlBackup = {
      enable = true;
      location = "/var/backups/postgresql";
      retentionDays = 14; # Keep backups for 2 weeks
    };

    # --- Django Systemd Service ---
    systemd.services.django-api = {
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" "postgresql.service" ]; # Wait for DB
      serviceConfig = {
        User = "django-user";
        Group = "django-user";
        WorkingDirectory = "/var/lib/django-api/myproject";
        ExecStart = "${pkgs.python312Packages.gunicorn}/bin/gunicorn \
          --workers 3 \
          --bind unix:/run/gunicorn.sock \
          myproject.wsgi:application";
        Restart = "on-failure";
        Environment = [
          "DATABASE_URL=postgres://django-user:1@localhost/bogbon"
          "SECRET_KEY=django-insecure-xter)8yl_*e(@&zybyny+cicwooq#g3l$mk#+qsjvyd1ldrmk4"
        ];
      };

      preStart = ''
        # Clone the repo
        if [ ! -d "/var/lib/django-api/myproject" ]; then
          git clone https://github.com/zawkindev/bogbon /var/lib/django-api/myproject
          chown -R django-user:django-user /var/lib/django-api/myproject
        fi
        cd /var/lib/django-api/myproject

        ${pkgs.python312}/bin/python manage.py migrate --noinput
        ${pkgs.python312}/bin/python manage.py collectstatic --noinput
      '';
    };

    # --- Nginx with HTTPS ---
    services.nginx = {
      enable = true;
      recommendedProxySettings = true;
      virtualHosts."bogbon.uz" = {
        forceSSL = true; # Redirect HTTP → HTTPS
        enableACME = true; # Auto SSL certs via Let's Encrypt
        locations."/" = {
          proxyPass = "http://unix:/run/gunicorn.sock";
          proxySetHeaders = {
            Host = "$host";
            X-Real-IP = "$remote_addr";
            X-Forwarded-Proto = "$scheme";
          };
        };
        locations."/static/" = {
          root = "/var/lib/django-api/myproject";
        };
      };
    };

    # --- Let's Encrypt SSL Certs ---
    security.acme = {
      acceptTerms = true;
      email = "zawkindev@gmail.com";
    };

    # --- User & Security ---
    users.users.django-user = {
      isSystemUser = true;
      group = "django-user";
      home = "/var/lib/django-api";
      createHome = true;
    };

    networking.firewall.allowedTCPPorts = [ 80 443 ]; # HTTP/HTTPS
  };
}
