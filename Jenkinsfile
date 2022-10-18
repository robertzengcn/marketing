pipeline {
    agent { label 'jenkins-agent' } 
    stages {
        stage('build docker composer') {
            steps {
                 sh 'docker build -t amigamarketing .' 
            }
        }
    }
}