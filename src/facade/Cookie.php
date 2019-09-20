<?php

namespace lumen\extra\facade;

use Lcobucci\JWT\Token;
use Illuminate\Support\Facades\Facade;

/**
 * Class Cookie
 * @method static bool set(string $name, $value, array $option = [])
 * @method static array|string|null get(string $name)
 * @method static bool delete(string $name)
 * @package lumen\extra\facade
 */
class Cookie extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return 'single-cookie';
    }
}
