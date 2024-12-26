{ lib, buildGoModule }:
buildGoModule {
  pname = "nix-converter";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-h6NQcwTI9SPWzgnIrQb5iOnSIPFjw1AJdlgyA+bmXW0=";

  meta = {
    description = "All-in-one converter from data format (JSON, YAML, etc.) to Nix and vice versa";
    homepage = "https://github.com/theobori/nix-converter";
    license = lib.licenses.mit;
    mainProgram = "nix-converter";
  };
}
