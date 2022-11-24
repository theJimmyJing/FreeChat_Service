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
	docker rmi "openim_${i}:$oldVersion"
	docker rmi "${ECR}/openim_${i}:$oldVersion"

	image="openim_${i}:$version"
	docker build -t $image . -f ./${i}.Dockerfile

	tag="${ECR}/openim_${i}:$version"
	aws ecr create-repository --repository-name "openim_${i}"
	docker tag ${image} ${tag}
	docker push ${tag}
	echo "push ${tag} success "

	rm -rf ./config.yaml
	cd ..
done



