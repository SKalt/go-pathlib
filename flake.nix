{
  description = "TODO";
  inputs = {
    flake-utils.url = "github:numtide/flake-utils"; # TODO: pin
    nixpkgs.url = "github:nixos/nixpkgs/nixos-23.11";
    # FIXME: eliminate reliance on GitHub
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = { self, flake-utils, nixpkgs, gomod2nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = (import nixpkgs) {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
      in {
        packages.default = pkgs.callPackage ./. { };
        devShell = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go # 1.21
          ];
          buildInputs = with pkgs; [
            # nix support
            nixpkgs-fmt
            nil

            # go development
            gopls
            # goreleaser # TODO: when we have a release

            # general development
            git
            bashInteractive
            lychee
            shellcheck
            fzf
            ripgrep

            # # node development
            # nodejs
            # nodePackages.pnpm
          ];
        };
      }
    );
}
