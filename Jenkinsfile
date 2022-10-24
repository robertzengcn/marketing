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
                 sh 'docker build --no-cache -t amigamarketing .' 
            }
        }
    }
}