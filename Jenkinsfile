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
    stage('Build Easyvpn cli for OSX') {
      steps {
        sh 'make init_osx'
      }
    }
    stage('Build Easyvpn cli for Linux') {
      steps {
        sh 'make init_linux'
      }
    }
    stage('Build Easyvpn cli for Windows') {
      steps {
        sh 'make init_windows'
      }
    }
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
  }
}

