{
  description = "Flake utils demo";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        formatter = pkgs.nixpkgs-fmt;
        packages = flake-utils.lib.flattenTree rec {
          pad-archiver = pkgs.callPackage ./nix { };
          dockerImage = pkgs.dockerTools.buildLayeredImage {
            name = "pad-archiver";
            tag = "latest";
            contents = [ pad-archiver ];
            extraCommands = ''
              mkdir -p repo
            '';
            config = {
              WorkingDir = "/repo";
              Entrypoint = [ "${pad-archiver}/bin/pad-archiver" ];
            };
          };
        };
        checks = {
          nix-format = pkgs.runCommand "nix-format" { buildInputs = [ formatter ]; } "nixpkgs-fmt --check ${./.} && touch $out";
        };
        defaultPackage = packages.pad-archiver;
        apps.pad-archiver = flake-utils.lib.mkApp { drv = packages.pad-archiver; };
        defaultApp = apps.pad-archiver;
        devShell = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            gopls
          ];
        };
      }
    );
}
