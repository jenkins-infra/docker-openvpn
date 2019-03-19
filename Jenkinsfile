pipeline {
  agent {
    label 'docker'
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
          sh 'make build.docker'
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
      parralel {
        stage('Build for OSX') {
          agent {
            label 'linux'
          }
          steps {
            sh 'make init_osx'
          }
        }
        stage('Build for Linux') {
          agent {
            label 'linux'
          }
          steps {
            sh 'make init_linux'
          }
        }
        stage('Build for Windows') {
          agent {
            label 'windows'
          }
          steps {
            sh 'make init_windows'
          }
        }
      }
    }
  }
}

