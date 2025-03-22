{
  description = "Dev environment with Nix";

  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      devShells = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.postgresql_15
              pkgs.python312
              pkgs.python312Packages.virtualenv
              pkgs.python312Packages.django
              pkgs.python312Packages.psycopg2
              pkgs.python312Packages.djangorestframework
            ];

            shellHook = ''
              export PGDATA="$PWD/.pgdata"
              export PATH="${pkgs.postgresql_15}/bin:$PATH"

              if [ ! -d "$PGDATA" ]; then
                echo "[INFO] Initializing PostgreSQL..."
                initdb --auth=md5 -D "$PGDATA"

                # WARNING: Password '1' is for development only!
                pg_ctl -D "$PGDATA" -l logfile start
                sleep 2
                psql -U postgres -c "ALTER USER postgres WITH PASSWORD '1';"
                createdb -U postgres postgres
                pg_ctl -D "$PGDATA" stop
              fi

              pg_ctl -D "$PGDATA" -l logfile start
              trap "pg_ctl -D '$PGDATA' stop" EXIT

              echo "[INFO] PostgreSQL ready. Connect with: psql -U postgres"
            '';
          };
        });
    };
}
