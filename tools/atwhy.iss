#include "environment.iss"

[Setup]
ChangesEnvironment=true
Compression=bzip/9
OutPutDir=installer
OutputBaseFilename=AtWhySetup
SourceDir=install-dir
UseSetupLdr=YES
PrivilegesRequired=admin
AppName=AtWhy
AppId=AtWhy
AppVersion=0.1.0.0
AppVerName=AtWhy version 0.1.0.0
AppMutex=AtWhy_Mutex
ChangesAssociations=YES
DefaultDirName={commonpf}\AtWhy
DefaultGroupName=AtWhy
DisableStartupPrompt=YES
;MessagesFile=C:\Inno Setup 3\Default.isl
AppCopyright=Tiffinger & Thiel GmbH © 2021
;BackColor=$FCF9DC
;BackColor2=$B05757
;windowVisible=YES
;WizardImageFile=C:\AtWhy\Setup\WizModernImage3.bmp
;WizardSmallImageFile=C:\AtWhy\Setup\WizModernSmallImage3.bmp
UserInfoPage=NO
DisableDirPage=NO
DisableReadyPage=NO
UsePreviousAppDir=YES
UninstallFilesDir={commonpf}\AtWhy\Uninstall Information
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
