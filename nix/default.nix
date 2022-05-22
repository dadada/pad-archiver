{ pkgs, lib, buildGoModule }:
buildGoModule {
  pname = "pad-archiver";
  version = "0.0.1";
  src = ../.;
  vendorSha256 = "sha256-e9I2mhSjLxMuCO9+g13XbDYI15Q879iG1AGZv6otuEA=";
  meta = with lib; {
    description = "Archives Etherpads with git";
    homepage = "https://git.fginfo.tu-bs.de/fginfo/pad-archiver";
    license = licenses.mit;
    maintainers = [ "y0067212" ];
  };
}
