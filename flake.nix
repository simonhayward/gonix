{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_25
            sqlite        # Provides SQLite 3.45+
            pkg-config    # Needed if building go-sqlite3 with CGO
          ];

          shellHook = ''
            echo "SQLite version: $(sqlite3 --version)"
            echo "Go version: $(go version)"
          '';
        };
      });
}