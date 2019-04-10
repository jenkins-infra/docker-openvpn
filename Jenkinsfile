pipeline {
  agent {
    label 'docker'
  }

  options {
    buildDiscarder logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')
    timeout(time: 1, unit: 'HOURS')
    timestamps()
  }

  triggers {
    pollSCM('H/15 * * * *')
  }

  stages {
    stage("Build Easyvpn"){
      failFast true
      parallel {
        stage('OSX') {
          steps {
            sh 'make init_osx'
          }
        }
        stage('Linux') {
          steps {
            sh 'make init_linux'
          }
        }
        stage('Windows') {
          steps {
            sh 'make init_windows'
          }
        }
      }
    }
    stage('Build OpenVPN Docker Image') {
      steps {
          sh 'make build.docker'
      }
    }
    stage('Publish OpenVPN Docker Image'){
      when {
        environment name: 'JENKINS_URL', value: 'https://trusted.ci.jenkins.io:1443/'
      }
      steps {
        script {
          infra.withDockerCredentials {
            sh 'make publish.docker'
          }
        }
      }
    }
  }
}

