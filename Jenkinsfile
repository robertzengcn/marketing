pipeline {
    agent { label 'jenkins-agent' } 

    stages {
       stage("provide config file"){
        steps{
      configFileProvider([configFile(fileId: 'amiga-marketing-conf', targetLocation: 'conf/app.conf')]) {   
             }
      }
      }
        stage('build docker composer') {
            steps {
                 sh 'docker image rm amigamarketing 2>&1 > /dev/null'
                 sh 'docker build -t amigamarketing .' 
            }
        }
    }
}