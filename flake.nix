{
  description = "Kmer counter";

  inputs = {
    mars-std.url = "github:mars-research/mars-std";
  };

  outputs = { self, mars-std, ... }: let
  # System types to support.
  supportedSystems = [ "x86_64-linux" ];

  in mars-std.lib.eachSystem supportedSystems (system: let
    pkgs = mars-std.legacyPackages.${system};
    in rec {
      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          act
          go_1_18
          golangci-lint
          nodePackages.serverless
        ];  
      };
    }
  );
}