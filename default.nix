{ lib, buildGoModule }:
buildGoModule {
  pname = "nix-converter";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-Ay1f9sk8RuJyOS7hl/lrscpxdlIgm9dMow/xTFoR+H4=";

  meta = {
    description = "All-in-one converter from data format (JSON, YAML, etc.) to Nix and vice versa";
    homepage = "https://github.com/theobori/nix-converter";
    license = lib.licenses.mit;
    mainProgram = "nix-converter";
  };
}
