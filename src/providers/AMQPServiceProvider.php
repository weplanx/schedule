<?php

namespace lumen\extra\providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use tidy\amqp\RabbitClient;

class AMQPServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     *
     * @return void
     */
    public function register()
    {
        $this->app->singleton('amqp', function (Application $app) {
            $config = $app->make('config')
                ->get('queue.connections.rabbitmq');

            return new RabbitClient($config);
        });
    }
}
