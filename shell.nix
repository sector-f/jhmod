with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "jhmod";
  buildInputs = [
    go_1_18
    shellcheck
  ];
}
