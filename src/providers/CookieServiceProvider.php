<?php

namespace lumen\extra\providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use lumen\extra\common\AuthFactory;
use lumen\extra\common\CookieFactory;
use Predis\Client;

class CookieServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        $this->app->singleton('single-cookie', function (Application $app) {
            $request = $app->make('request');
            $config = $app->make('config')
                ->get('cookie');

            return new CookieFactory($request, $config);
        });
    }
}
