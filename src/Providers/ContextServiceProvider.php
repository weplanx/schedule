<?php

namespace Lumen\Extra\Providers;

use Laravel\Lumen\Application;
use Illuminate\Support\ServiceProvider;
use Lumen\Extra\Common\ContextFactory;
use Lumen\Extra\Common\TokenFactory;

final class ContextServiceProvider extends ServiceProvider
{
    /**
     * Register the service provider.
     * @return void
     */
    public function register()
    {
        $this->app->singleton('context', function (Application $app) {
            return new ContextFactory(
            );
        });
    }
}
