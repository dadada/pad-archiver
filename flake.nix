{
  description = "Flake utils demo";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      rec {
        packages = flake-utils.lib.flattenTree {
          pad-archiver = pkgs.callPackage ./nix { pkgs = pkgs; };
          gitAndTools = pkgs.gitAndTools;
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
