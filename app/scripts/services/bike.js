angular.module('bikeService', ['ngResource']).factory('Bike', function($resource){
		console.log($resource);
		return $resource('http://192.168.198.130:8080/bikes/:id',{},{
			get:{isArray:true},
			put:{method:'PUT', isArray:false},
			update:{method: 'PUT'},

		});
	});