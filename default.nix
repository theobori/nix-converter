{ lib, buildGoModule }:
buildGoModule {
  pname = "data2nix";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-4q0K+3MjYOkg3hDdAePLkEDf3lcrP6ng48ZAemMdv7g=";

  meta = {
    description = "Convet configuration languages (JSON, YAML, etc..) to Nix";
    homepage = "https://github.com/theobori/data2nix";
    license = lib.licenses.mit;
    mainProgram = "data2nix";
  };
}
