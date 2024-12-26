{ ... }:
{
  projectRootFile = "flake.nix";
  programs.black.enable = true;
  programs.gofmt.enable = true;
}
