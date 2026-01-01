{ ... }:
{
  projectRootFile = "flake.nix";
  programs.nixfmt.enable = true;
  programs.gofmt.enable = true;
  programs.actionlint.enable = true;
}
