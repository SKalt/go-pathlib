{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
  , buildGoApplication ? pkgs.buildGoApplication
}:
let
  filters = import ./nix/filter_filesets.nix { inherit (pkgs) lib; };
in
  pkgs.buildGoApplication {
    pname = "TODO";
    version = "0.0.0"; # TODO: update mechanism
    pwd = ./.;
    src = filters.go ./.;
    modules = ./gomod2nix.toml;
  }
