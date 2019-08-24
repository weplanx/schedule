<?php

namespace lumen\extra\common;

use Illuminate\Http\Request;

final class CookieFactory
{
    /**
     * App Request
     * @var Request $request
     */
    private $request;

    /**
     * Cookie Config
     * @var array $config
     */
    private $config;

    /**
     * CookieFactory constructor.
     * @param array $config
     */
    public function __construct(Request $request, array $config)
    {
        $this->request = $request;
        $this->config = $config;
    }

    /**
     * Set Cookie
     * @param string $name cookie key
     * @param mixed $value cookie value
     * @param array $option cookie option
     * @return bool
     */
    public function set(string $name, $value, array $option = [])
    {
        $args = array_merge(
            $this->config,
            $option
        );
        return setcookie(
            $name,
            $value,
            $args['expire'],
            $args['path'],
            $args['domain'],
            $args['secure'],
            $args['httponly']
        );
    }

    /**
     * Get Cookie
     * @param string $name cookie key
     * @return array|string|null
     */
    public function get(string $name)
    {
        return $this->request->cookie($name);
    }

    /**
     * Delete Cookie
     * @param string $name cookie key
     * @return bool
     */
    public function delete(string $name)
    {
        return $this->set($name, null);
    }
}
