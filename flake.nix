{
  description = "Generate templated files quickly from your terminal.";

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let 
	    pkgs = nixpkgs.legacyPackages.${system}; 
      in
      {
        packages = rec { 
            plate = pkgs.buildGoModule rec {
                name = "plate";
                src = ./.;
                vendorHash = "sha256-vanKL5s+szW0hduUXGnJNUlyu8wZ2HsBVklIUb/+DLY=";
            };
            default = plate;
        };
        devShells.default = pkgs.mkShell {
          packages = [ pkgs.go ];
        };
      }
    );
}
