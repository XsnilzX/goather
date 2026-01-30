{
  description = "Weather widget for Quickshell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    # Dev Shell (zum Entwickeln)
    devShells.${system}.default = pkgs.mkShell {
      packages = with pkgs; [
        go
        gopls # Go Language Server
        gotools # goimports, etc.
      ];
    };

    # Package (zum Nutzen)
    packages.${system} = {
      goather = pkgs.buildGoModule {
        pname = "goather";
        version = "1.0.0";
        src = ./.;

        vendorHash = null; # Ändern falls go.sum existiert

        ldflags = ["-s" "-w"];

        meta = {
          description = "Weather widget for Quickshell/Hyprland";
          mainProgram = "goather";
        };
      };

      default = self.packages.${system}.goather;
    };

    # Für NixOS/Home-Manager
    nixosModules.default = {
      config,
      pkgs,
      ...
    }: {
      environment.systemPackages = [
        self.packages.${system}.goather
      ];
    };

    homeManagerModules.default = {
      config,
      pkgs,
      ...
    }: {
      home.packages = [
        self.packages.${system}.goather
      ];
    };
  };
}
