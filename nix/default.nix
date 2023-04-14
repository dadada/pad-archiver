{ pkgs, lib, buildGoModule }:
buildGoModule {
  pname = "pad-archiver";
  version = "0.0.1";
  src = ../.;
  vendorSha256 = "sha256-vBni3j3o0P13PJg/Ab1ux9zSVr05Iha/sb8dVTX4G0g=";
  meta = with lib; {
    description = "Archives Etherpads with git";
    homepage = "https://git.fginfo.tu-bs.de/fginfo/pad-archiver";
    license = licenses.mit;
    maintainers = [ "y0067212" ];
  };
}
