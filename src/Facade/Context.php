<?php

namespace Lumen\Extra\Facade;

use Illuminate\Support\Facades\Facade;
use Lumen\Extra\Contracts\ContextInterface;

/**
 * Class Context
 * @package Lumen\Extra\Facade
 * @method static void set(string $abstract, $value)
 * @method static mixed get($abstract)
 */
final class Context extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return ContextInterface::class;
    }
}
