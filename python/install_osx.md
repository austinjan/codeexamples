# Install python in macos
In macos, os install python2 and python3 by default. But default python will not update to latest version. So we need to install python3 manually.


## Pre-requisites
- Homebrew:  `brew update` to update homebrew to latest version


## Install python3
```
brew install python
```

Setup aliases for python3 and pip3:  
```
echo "alias python=python3" >> ~/.zshrc
echo "alias pip=pip3" >> ~/.zshrc
```

## vs-code (optional)
Extensions:
- Python

Optional Extensions:
- vscode-icons