<?php

namespace Lumen\Extra\Facade;

use Illuminate\Support\Facades\Facade;

/**
 * Class Token
 * @package Lumen\Extra\Facade
 * @method static \Lcobucci\JWT\Token|false create(string $scene, string $jti, string $ack, array $symbol = [])
 * @method static \Lcobucci\JWT\Token get(string $tokenString)
 * @method static \stdClass verify(string $scene, string $tokenString)
 */
final class Token extends Facade
{
    /**
     * Get the registered name of the component.
     *
     * @return string
     */
    protected static function getFacadeAccessor()
    {
        return 'token';
    }
}
