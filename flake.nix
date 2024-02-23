{
  description = "Generate templated files quickly from your terminal.";

  outputs = { nixpkgs, ... }: let
    forAllSystems = function:
      nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
      ] (system: function nixpkgs.legacyPackages.${system});
  in rec {
    packages = forAllSystems (pkgs: rec {
      default = plate;
      plate = pkgs.buildGoModule rec {
        pname = "plate";
        version = "v0.0.1";
        src = ./.;
        vendorHash = "sha256-JFrzxduL0Wr3+CGfAmJbAcaCWRP/vLF6nQWds2aamtw=";
        ldflags = ["-X 'main.Version=${version}'"];
      };
    });
    devShells= forAllSystems (pkgs: {
      default = pkgs.mkShell {
        packages = [pkgs.go pkgs.gopls];
      };
    });
  };
}
