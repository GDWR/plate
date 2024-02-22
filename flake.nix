{
  description = "Template files in CLI for quick creation.";

  outputs = { self, nixpkgs, flake-utils }: 
    flake-utils.lib.eachDefaultSystem (system:
      let 
	    pkgs = nixpkgs.legacyPackages.${system}; 
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [ pkgs.go ];
        };
      }
    );
}
