with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "jh_extract";
  buildInputs = [
    go_1_18
    shellcheck
  ];
}
