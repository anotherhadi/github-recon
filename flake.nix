{
  description =
    "GH-Recon: Fetches and aggregates public OSINT data for a GitHub user, leveraging Go and the GitHub API.";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = f:
        nixpkgs.lib.genAttrs supportedSystems
        (system: f system (import nixpkgs { inherit system; }));

      pname = "gh-recon";
      version = "0.1.0";

      ldflags = [ "-s" "-w" ];

    in {
      packages = forAllSystems (system: pkgs: {
        "${pname}" = pkgs.buildGoModule {
          inherit pname version ldflags;

          src = ./.;

          vendorHash = "sha256-CPk8B8FKEoN8qff6WV/iBf0eVjTBMVfJQvlVcti6dfM=";

          meta = with pkgs.lib; {
            description =
              "Fetches and aggregates public OSINT data for a GitHub user.";
            homepage = "https://github.com/anotherhadi/gh-recon";
            platforms = platforms.unix;
          };
        };
      });

      defaultPackage =
        forAllSystems (system: pkgs: self.packages.${system}.${pname});
    };
}
