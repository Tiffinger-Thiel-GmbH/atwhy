#include "environment.iss"

[Setup]
ChangesEnvironment=true
Compression=bzip/9
OutPutDir=installer
OutputBaseFilename=atwhySetup
SourceDir=install-dir
UseSetupLdr=YES
PrivilegesRequired=admin
AppName=atwhy
AppId=atwhy
AppVersion=0.1.0.0
AppVerName=atwhy version 0.1.0.0
AppMutex=atwhy_Mutex
ChangesAssociations=YES
DefaultDirName={commonpf}\atwhy
DefaultGroupName=atwhy
DisableStartupPrompt=YES
;MessagesFile=C:\Inno Setup 3\Default.isl
AppCopyright=Tiffinger & Thiel GmbH © 2021
;BackColor=$FCF9DC
;BackColor2=$B05757
;windowVisible=YES
;WizardImageFile=C:\atwhy\Setup\WizModernImage3.bmp
;WizardSmallImageFile=C:\atwhy\Setup\WizModernSmallImage3.bmp
UserInfoPage=NO
DisableDirPage=NO
DisableReadyPage=NO
UsePreviousAppDir=YES
UninstallFilesDir={commonpf}\atwhy\Uninstall Information
ShowTasksTreeLines=YES

; More options in setup section as well as other sections like Files, Components, Tasks...

[Files]
Source: "atwhy.exe"; DestDir: "{app}"

[Tasks]
Name: envPath; Description: "Add to PATH variable"

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
    if (CurStep = ssPostInstall) and WizardIsTaskSelected('envPath')
     then EnvAddPath(ExpandConstant('{app}'));
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
    if CurUninstallStep = usPostUninstall
    then EnvRemovePath(ExpandConstant('{app}'));
end;
