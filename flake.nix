{
  description = "An example Go package";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils}:
    # Generate a user-friendly version number.
    let version = builtins.substring 0 8 self.lastModifiedDate;
    in {
      overlay = final: prev:
        let pkgs = nixpkgs.legacyPackages.${prev.system};
        in rec {
          gonix = pkgs.buildGo120Module rec {
            pname = "gonix";
            inherit version;
            src = pkgs.nix-gitignore.gitignoreSource [ ] ./.;

            vendorSha256 = "sha256-8SfbzJ/uER7LS9EhdhgG5K18WLC2AbJAPLozFyZ5mJo=";
          };
        };
    }
    //flake-utils.lib.eachDefaultSystem
    (system:
    let pkgs = import nixpkgs {
        overlays = [ self.overlay ];
        inherit system;
      };
    in rec {
      # `nix develop`
      devShell = pkgs.mkShell { buildInputs = with pkgs; [ go_1_20 gopls goimports go-tools ]; };

      # `nix build`
      packages = with pkgs; {
        inherit gonix;
      };

      defaultPackage = pkgs.gonix;

      # `nix run`
      apps.gonix = flake-utils.lib.mkApp {
        drv = packages.gonix;
      };
      defaultApp = apps.gonix;

      overlays.default = self.overlay;
    });
}