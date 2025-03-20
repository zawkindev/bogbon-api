{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  name = "django-env";

  buildInputs = [
    pkgs.python312
    pkgs.python312Packages.pip
    pkgs.python312Packages.virtualenv
    pkgs.python312Packages.django
    pkgs.python312Packages.psycopg2
    pkgs.python312Packages.djangorestframework
  ];

  shellHook = ''
    echo "🔧 Setting up Python virtual environment..."

    if [ ! -d .venv ]; then
      python -m venv .venv
    fi

    source .venv/bin/activate

    echo "✅ Virtual environment activated."
  '';
}

