#include "environment.iss"

[Setup]
ChangesEnvironment=true
Compression=bzip/9
OutPutDir=installer
OutputBaseFilename=CrazyDocSetup
SourceDir=install-dir
UseSetupLdr=YES
PrivilegesRequired=admin
AppName=CrazyDoc
AppId=CrazyDoc
AppVersion=0.1.0.0
AppVerName=CrazyDoc version 0.1.0.0
AppMutex=CrazyDoc_Mutex
ChangesAssociations=YES
DefaultDirName={commonpf}\CrazyDoc
DefaultGroupName=CrazyDoc
DisableStartupPrompt=YES
;MessagesFile=C:\Inno Setup 3\Default.isl
AppCopyright=Tiffinger & Thiel GmbH © 2021
;BackColor=$FCF9DC
;BackColor2=$B05757
;windowVisible=YES
;WizardImageFile=C:\CrazyDoc\Setup\WizModernImage3.bmp
;WizardSmallImageFile=C:\CrazyDoc\Setup\WizModernSmallImage3.bmp
UserInfoPage=NO
DisableDirPage=NO
DisableReadyPage=NO
UsePreviousAppDir=YES
UninstallFilesDir={commonpf}\CrazyDoc\Uninstall Information
ShowTasksTreeLines=YES

; More options in setup section as well as other sections like Files, Components, Tasks...

[Files]
Source: "crazydoc.exe"; DestDir: "{app}"

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
    if CurStep = ssPostInstall
     then EnvAddPath(ExpandConstant('{app}'));
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
    if CurUninstallStep = usPostUninstall
    then EnvRemovePath(ExpandConstant('{app}'));
end;
