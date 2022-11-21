#/bin/sh
source ./path_info.cfg

git pull
cd ../script/; ./build_all_service.sh
cd ../deploy_k8s/

for i in  ${service[*]}
do
  cp ../config/config.yaml ./${i}/
  mv ../bin/open_im_${i} ./${i}/
done

echo "move success"

aws ecr get-login-password | docker login --username AWS --password-stdin ${ECR}

echo "start to build images"
for i in ${service[*]}
do
	echo "start to build images" $i
	cd $i
	docker rmi "${ECR}/${i}:$oldVersion"
	image="${ECR}/${i}:$version"
	docker build -t $image . -f ./${i}.Dockerfile
	echo "build ${dockerfile} success"
	echo "clean temp success"
	docker push $image
	echo "push ${image} success "
	rm -rf ./config.yaml
	cd ..
done



