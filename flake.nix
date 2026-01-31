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
        packages.default = pkgs.buildGoModule rec{
          pname = "gonix";
          version = "0.1.0";
          src = ./.;

          # You must update this hash whenever go.mod changes.
          vendorHash = "sha256-UTkp3qXSpq/hljlAh4CWMhg4T0r7yJwDR/CPWqhtNe4=";

          preBuild = ''
            export CGO_ENABLED=1
          '';

          # self.lastModified is the Unix timestamp of the last commit
          # builtins.toString converts it for the ldflags
          ldflags = [
            "-X main.Version=${version}"
            "-X main.Commit=${self.rev or "dirty"}"
            "-X main.BuildTime=${builtins.toString self.lastModified}"
          ];

          buildInputs = [ pkgs.sqlite ];
        };

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
