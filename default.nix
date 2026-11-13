{ lib, buildGoModule }:
buildGoModule {
  pname = "nix-converter";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-SvJThxyduaBkN965suZzbsxN1LiRv79TS2JH7i78svI=";

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "All-in-one converter from data format (JSON, YAML, etc.) to Nix and vice versa";
    homepage = "https://github.com/theobori/nix-converter";
    license = lib.licenses.mit;
    mainProgram = "nix-converter";
  };
}
