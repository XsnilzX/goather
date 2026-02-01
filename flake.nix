{
  description = "Weather widget for Quickshell written in Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
      version = "1.0.2";
    in {
      # Package
      packages = {
        goather = pkgs.buildGoModule {
          pname = "goather";
          inherit version;

          src = pkgs.lib.cleanSourceWith {
            src = ./.;
            filter = path: type: let
              baseName = baseNameOf path;
            in
              !(pkgs.lib.hasSuffix ".md" baseName)
              && !(baseName == "flake.nix")
              && !(baseName == "flake.lock")
              && !(pkgs.lib.hasPrefix "." baseName);
          };

          vendorHash = null;

          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];

          # Build-Tags für standalone binary
          tags = ["netgo" "osusergo"];

          meta = with pkgs.lib; {
            description = "Weather widget for Quickshell/Hyprland";
            homepage = "https://github.com/XsnilzX/goather";
            license = licenses.mit;
            maintainers = [maintainers.XsnilzX]; # Optional: maintainers.DEIN_NAME
            mainProgram = "goather";
            platforms = platforms.linux;
          };
        };

        default = self.packages.${system}.goather;
      };

      # Dev Shell
      devShells.default = pkgs.mkShell {
        inputsFrom = [self.packages.${system}.goather];

        packages = with pkgs; [
          go
          gopls
          gotools
          go-tools
          delve
          golangci-lint
        ];

        shellHook = ''
          echo "🌤️  Goather Development Environment"
          echo "Go version: $(go version | awk '{print $3}')"
          echo ""
          echo "Commands:"
          echo "  go run .              - Run application"
          echo "  go test ./...         - Run tests"
          echo "  golangci-lint run     - Lint code"
          echo "  nix build             - Build package"
        '';
      };

      # Apps (direkt ausführbar mit: nix run)
      apps.default = {
        type = "app";
        program = "${self.packages.${system}.goather}/bin/goather";
      };
    })
    # System-unabhängige Outputs
    // {
      # Overlay für einfache Integration
      overlays.default = final: prev: {
        goather = self.packages.${final.system}.goather;
      };

      # NixOS Module
      nixosModules.default = {
        config,
        lib,
        pkgs,
        ...
      }: let
        cfg = config.programs.goather;
      in {
        options.programs.goather = {
          enable = lib.mkEnableOption "Goather weather widget";

          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.goather;
            description = "The goather package to use";
          };
        };

        config = lib.mkIf cfg.enable {
          environment.systemPackages = [cfg.package];
        };
      };

      # Home Manager Module
      homeManagerModules.default = {
        config,
        lib,
        pkgs,
        ...
      }: let
        cfg = config.programs.goather;
      in {
        options.programs.goather = {
          enable = lib.mkEnableOption "Goather weather widget";

          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.goather;
            description = "The goather package to use";
          };
        };

        config = lib.mkIf cfg.enable {
          home.packages = [cfg.package];
        };
      };
    };
}
