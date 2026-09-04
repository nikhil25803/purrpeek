{
  description = "Your system information, from a cat's perspective";

  inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-26.05-darwin";
  };

  outputs = { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = nixpkgs.lib.genAttrs systems;

    in {
      packages = forAllSystems (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in {
          default = pkgs.buildGoModule {
            pname = "purrpeek";
            version = "0.1.0";

            src = ./.;

            vendorHash = null;

            meta = {
              description = "Your system information, from a cat's perspective";
              homepage = "https://github.com/nikhil25803/purrpeek";
              license = nixpkgs.lib.licenses.mit;
              mainProgram = "purrpeek";
            };
          };
        }
      );
    };
}