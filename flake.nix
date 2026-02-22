{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachSystem  [ "x86_64-linux" ] (system:
      let
        pkgs = import nixpkgs { inherit system; };

        appVersion = 
          let 
            envV = builtins.getEnv "APP_VERSION"; 
          in
            if envV != "" then envV 
            else if (self ? shortRev) then self.shortRev 
            else "dev";

        gonix = pkgs.buildGo125Module {
          pname = "gonix";
          version = appVersion;
          src = ./.;

          vendorHash = "sha256-UTkp3qXSpq/hljlAh4CWMhg4T0r7yJwDR/CPWqhtNe4="; # update this hash for go.mod changes.

          # static build for Distroless compatibility
          preBuild = ''
            export CGO_ENABLED=0
          '';

          ldflags = [
            "-X main.Version=${appVersion}"
            "-X main.Commit=${self.rev or "dirty"}"
            "-X main.BuildTime=${builtins.toString self.lastModified}"
          ];

          buildInputs = [ pkgs.sqlite ];
        };
      in
      {
        devShells.default = pkgs.mkShell {
            buildInputs = with pkgs; [
              bashInteractive
              go_1_25
              sqlite
            ];

           shellHook = ''
             echo "SQLite version: $(sqlite3 --version)"
             echo "Go version: $(go version)"
           '';
        };

        packages = {
          default = gonix;
          dockerImage = pkgs.dockerTools.buildLayeredImage {
            name = "registry.fly.io/gonix-deploy";
            tag = "${appVersion}";
            contents = [ gonix ];

            fakeRootCommands = ''
              mkdir -p tmp
              chmod 1777 tmp
            '';

            config = {
              Cmd = [ "/bin/gonix" ];
              Env = [ "TMPDIR=/tmp" "HOME=/home/nonroot" ];
              WorkingDir = "/home/nonroot";
            };
          };
        };
      });
}
