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

        # nix run 'nixpkgs#nix-prefetch-docker' -- --image-name gcr.io/distroless/static-debian12 --image-tag nonroot-amd64
        distrolessBase = pkgs.dockerTools.pullImage {
          imageName = "gcr.io/distroless/static-debian12";
          imageDigest = "sha256:5074667eecabac8ac5c5d395100a153a7b4e8426181cca36181cd019530f00c8";
          hash = "sha256-k00a6y/QB4UC9PM9fiQy0nK9trNdXo3KOP8nF0B3+iE=";
          finalImageName = "gcr.io/distroless/static-debian12";
          finalImageTag = "nonroot-amd64";
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
            fromImage = distrolessBase;

            config = {
              Cmd = [ "/bin/gonix" ];
              Env = [ "HOME=/home/nonroot" ];
              WorkingDir = "/home/nonroot";
            };
          };
        };
      });
}
