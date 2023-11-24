codesign --force -s "Developer ID Application: Gergely Svigruha (K74N482UPR)" -v Medville.app --deep --strict --options=runtime --timestamp
/usr/bin/ditto -c -k --keepParent Medville.app Medville.zip
xcrun notarytool submit Medville.zip --apple-id sgergo88@gmail.com --team-id K74N482UPR --wait --password $APP_SPEC_PWD
