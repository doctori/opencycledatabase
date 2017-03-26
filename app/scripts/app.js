'use strict';

/**
 * @ngdoc overview
 * @name openBicycleDatabaseApp
 * @description
 * # openBicycleDatabaseApp
 *
 * Main module of the application.
 */
angular
  .module('openBicycleDatabaseApp', [
    'bikeService',
    'componentService',
    'ngAnimate',
    'ngCookies',
    'ngMaterial',
    'ngResource',
    'ngRoute',
    'ngSanitize',
 //   'ngTouch',
    'gettext',
    'xeditable',
    'ngFileUpload'
  ])
  .run(function(gettextCatalog){
    gettextCatalog.setCurrentLanguage('fr_FR');
    gettextCatalog.debug = true;
  })
  .config(function ($routeProvider,$locationProvider,$provide) {
    $provide.decorator('$sniffer', function($delegate) {
      $delegate.history = false;
      return $delegate;
    });
    $locationProvider
      .html5Mode(true)
      .hashPrefix('!');
    $routeProvider
      .when('/', {
        templateUrl: 'views/bike.html',
        controller: 'BikeListCtrl',
        controllerAs: 'bike'
      })
      .when('/bike/:bikeID',{
        templateUrl: 'views/bike-detail.html',
        controller: 'BikeDetailCtrl',
        controllerAs: 'bike'
      })
      .when('/new-bike',{
        templateUrl: 'views/new-bike.html',
        controller: 'NewBikeCtrl',
        controllerAs: 'bike'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl',
        controllerAs: 'about'
      })
      .when('/component/:componentID', {
        templateUrl: 'views/component-detail.html',
        controller: 'ComponentDetailCtrl',
        controllerAs: 'component'
      })
      .when('/components', {
        templateUrl: 'views/components.html',
        controller: 'ComponentListCtrl',
        controllerAs: 'component'
      })

      .otherwise({
        redirectTo: '/'
      });
  });