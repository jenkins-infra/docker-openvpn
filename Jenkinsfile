pipeline {
  agent {
    docker
  }

  options {
    buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')
  }

  triggers {
    cron 'H/15 * * * *'
  }
  stages {
    stage('Build OpenVPN Docker Image') {
      steps {
          sh 'make build'
        }
      }
    }
    stage('Publish OpenVPN Docker Image'){
      when {
        branch 'master'
      }
      steps {
        make publish.docker
      }
    }
    stage('Build Easyvpn Cli'){
      steps{
        parralel (
          "OSX": {
            make init_osx
          }
          "Windows": {
            make init_windows
          }
          "Linux": {
            make init_linux
          }
        )
      }
    }
  }
}

