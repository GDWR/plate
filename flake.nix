{
  description = "Generate templated files quickly from your terminal.";
  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        formatter = pkgs.alejandra;
        packages = rec {
          plate = pkgs.buildGoModule rec {
            pname = "plate";
            version = "v0.0.1";
            src = ./.;
            vendorHash = "sha256-UgLotBPiykyV4mksUuX3697et2iTJVOYroL8C5Wd9qg=";
            ldflags = ["-X 'main.Version=${version}'"];
          };
          default = plate;
        };
        devShells.default = pkgs.mkShell {
          packages = [pkgs.go pkgs.gopls];
        };
      }
    );
}
