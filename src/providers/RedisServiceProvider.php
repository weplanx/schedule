<?php

namespace lumen\extra\providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use lumen\extra\common\AuthFactory;
use Predis\Client;

class RedisServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        $this->app->singleton('redis', function (Application $app) {
            $config = $app->make('config')
                ->get('database.redis');

            return new Client($config);
        });
    }
}
