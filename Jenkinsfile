pipeline {
    agent { label 'local-agent' } 
    stages {
        stage('build docker composer') {
            steps {
                 sh 'docker build -t amigamarketing .' 
            }
        }
    }
}