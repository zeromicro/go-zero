{
  description = "goctl's nix flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, utils, ... }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        defaultPkg = pkgs.callPackage ./default.nix {
          inherit pkgs;
        };
      in {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            # Development tools
            go

            # Aux tools
            goreleaser
          ];
        };

        packages = {
          default = defaultPkg;
        };
    });
}
